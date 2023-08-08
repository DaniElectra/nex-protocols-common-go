package matchmake_extension

import (
	nex "github.com/PretendoNetwork/nex-go"
	match_making_types "github.com/PretendoNetwork/nex-protocols-go/match-making/types"
	matchmake_extension "github.com/PretendoNetwork/nex-protocols-go/matchmake-extension"
	"github.com/PretendoNetwork/plogger-go"
)

var commonMatchmakeExtensionProtocol *CommonMatchmakeExtensionProtocol
var logger = plogger.NewLogger()

type CommonMatchmakeExtensionProtocol struct {
	*matchmake_extension.Protocol
	server *nex.Server

	cleanupSearchMatchmakeSessionHandler         func(matchmakeSession *match_making_types.MatchmakeSession)
	cleanupMatchmakeSessionSearchCriteriaHandler func(lstSearchCriteria []*match_making_types.MatchmakeSessionSearchCriteria)
}

// CleanupSearchMatchmakeSession sets the CleanupSearchMatchmakeSession handler function
func (commonMatchmakeExtensionProtocol *CommonMatchmakeExtensionProtocol) CleanupSearchMatchmakeSession(handler func(matchmakeSession *match_making_types.MatchmakeSession)) {
	commonMatchmakeExtensionProtocol.cleanupSearchMatchmakeSessionHandler = handler
}

// CleanupMatchmakeSessionSearchCriteria sets the CleanupMatchmakeSessionSearchCriteria handler function
func (commonMatchmakeExtensionProtocol *CommonMatchmakeExtensionProtocol) CleanupMatchmakeSessionSearchCriteria(handler func(lstSearchCriteria []*match_making_types.MatchmakeSessionSearchCriteria)) {
	commonMatchmakeExtensionProtocol.cleanupMatchmakeSessionSearchCriteriaHandler = handler
}

// NewCommonMatchmakeExtensionProtocol returns a new CommonMatchmakeExtensionProtocol
func NewCommonMatchmakeExtensionProtocol(server *nex.Server) *CommonMatchmakeExtensionProtocol {
	MatchmakeExtensionProtocol := matchmake_extension.NewProtocol(server)
	commonMatchmakeExtensionProtocol = &CommonMatchmakeExtensionProtocol{Protocol: MatchmakeExtensionProtocol, server: server}

	MatchmakeExtensionProtocol.AutoMatchmakePostpone(autoMatchmake_Postpone)
	MatchmakeExtensionProtocol.AutoMatchmakeWithSearchCriteriaPostpone(autoMatchmakeWithSearchCriteria_Postpone)
	MatchmakeExtensionProtocol.OpenParticipation(openParticipation)

	return commonMatchmakeExtensionProtocol
}