package bootstrap

import (
	"errors"
	"sync"

	"github.com/webrtc-demo-go/types"
)

// WsLinkMgr Управляйте отношениями сопоставления между идентификатором пользовательского агента внутреннего браузера и соединением WebSocket, добавляйте, удаляйте, проверяйте
// Режим связи mqtt является асинхронным, и должны быть отношения сопоставления, чтобы после получения сообщения mqtt он знал, какому клиенту WebSocket отправить ответ.
type WsLinkMgr struct {
	rwMutex sync.RWMutex

	session2Agent map[string]string

	wsLink map[string]*types.WsMessage // agent id -> WebSocket连接
}

var wsLinkMgr *WsLinkMgr

func init() {
	wsLinkMgr = &WsLinkMgr{
		session2Agent: make(map[string]string),

		wsLink: make(map[string]*types.WsMessage),
	}
}

// AddLink Увеличьте соединение WebSocket, связанное с идентификатором пользовательского агента браузера.
func AddLink(agentID, sessionID string, msg *types.WsMessage) {
	wsLinkMgr.rwMutex.Lock()
	defer wsLinkMgr.rwMutex.Unlock()

	// Идентификатор агента соответствует веб-странице, и веб-страница может иметь только один сеанс одновременно, и сеанс, связанный до того, как идентификатор агента будет очищен.
	for session, agent := range wsLinkMgr.session2Agent {
		if agent == agentID {
			delete(wsLinkMgr.session2Agent, session)
		}
	}

	wsLinkMgr.session2Agent[sessionID] = agentID

	wsLinkMgr.wsLink[agentID] = msg
}

// GetLink Запросите соединение WebSocket, связанное с идентификатором пользовательского агента браузера.
func GetLink(sessionID string) (link *types.WsMessage, err error) {
	wsLinkMgr.rwMutex.RLock()
	defer wsLinkMgr.rwMutex.RUnlock()

	agentID, ok := wsLinkMgr.session2Agent[sessionID]
	if !ok {
		return nil, errors.New("get agent fail")
	}

	link, ok = wsLinkMgr.wsLink[agentID]
	if !ok {
		return nil, errors.New("getLink fail")
	}

	return
}

func GetLinkByAgent(agentID string) (link *types.WsMessage, err error) {
	wsLinkMgr.rwMutex.RLock()
	defer wsLinkMgr.rwMutex.RUnlock()

	ok := false

	link, ok = wsLinkMgr.wsLink[agentID]
	if !ok {
		return nil, errors.New("getLink fail")
	}

	return
}

// RemoveLink Удалите соединение WebSocket, связанное с идентификатором пользовательского агента браузера.
func RemoveLink(sessionID string) {
	wsLinkMgr.rwMutex.Lock()
	defer wsLinkMgr.rwMutex.Unlock()

	agentID, ok := wsLinkMgr.session2Agent[sessionID]
	if !ok {
		return
	}

	delete(wsLinkMgr.session2Agent, sessionID)

	_, ok = wsLinkMgr.wsLink[agentID]
	if ok {
		delete(wsLinkMgr.wsLink, agentID)
	}
}

// RemoveLinkByConnLost Когда веб-страница браузера отключает соединение WebSocket, очистите соответствующие записи.
func RemoveLinkByConnLost(agentID string) {
	wsLinkMgr.rwMutex.Lock()
	defer wsLinkMgr.rwMutex.Unlock()

	for session, agent := range wsLinkMgr.session2Agent {
		if agent == agentID {
			delete(wsLinkMgr.session2Agent, session)
		}
	}

	if _, ok := wsLinkMgr.wsLink[agentID]; ok {
		delete(wsLinkMgr.wsLink, agentID)
	}
}
