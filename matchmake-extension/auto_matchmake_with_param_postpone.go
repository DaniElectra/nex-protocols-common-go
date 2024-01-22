package matchmake_extension

import (
	"github.com/PretendoNetwork/nex-go"
	"github.com/PretendoNetwork/nex-go/types"
	common_globals "github.com/PretendoNetwork/nex-protocols-common-go/globals"
	match_making_types "github.com/PretendoNetwork/nex-protocols-go/match-making/types"
	matchmake_extension "github.com/PretendoNetwork/nex-protocols-go/matchmake-extension"
)

func autoMatchmakeWithParam_Postpone(err error, packet nex.PacketInterface, callID uint32, autoMatchmakeParam *match_making_types.AutoMatchmakeParam) (*nex.RMCMessage, uint32) {
	if err != nil {
		common_globals.Logger.Error(err.Error())
		return nil, nex.Errors.Core.InvalidArgument
	}

	// TODO - This assumes a PRUDP connection. Refactor to support HPP
	connection := packet.Sender().(*nex.PRUDPConnection)
	endpoint := connection.Endpoint
	server := endpoint.Server

	// * A client may disconnect from a session without leaving reliably,
	// * so let's make sure the client is removed from the session
	common_globals.RemoveConnectionFromAllSessions(connection)

	matchmakeSession := autoMatchmakeParam.SourceMatchmakeSession

	sessions := common_globals.FindSessionsByMatchmakeSessionSearchCriterias(connection.PID(), autoMatchmakeParam.LstSearchCriteria.Slice(), commonProtocol.GameSpecificMatchmakeSessionSearchCriteriaChecks)
	var session *common_globals.CommonMatchmakeSession

	if len(sessions) == 0 {
		var errCode uint32
		session, err, errCode = common_globals.CreateSessionByMatchmakeSession(matchmakeSession, nil, connection.PID())
		if err != nil {
			common_globals.Logger.Error(err.Error())
			return nil, errCode
		}
	} else {
		session = sessions[0]
	}

	err, errCode := common_globals.AddPlayersToSession(session, []uint32{connection.ID}, connection, "")
	if err != nil {
		common_globals.Logger.Error(err.Error())
		return nil, errCode
	}

	matchmakeDataHolder := types.NewAnyDataHolder()

	matchmakeDataHolder.TypeName = types.NewString("MatchmakeSession")
	matchmakeDataHolder.ObjectData = session.GameMatchmakeSession.Copy()

	rmcResponseStream := nex.NewByteStreamOut(server)

	session.GameMatchmakeSession.WriteTo(rmcResponseStream)

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCSuccess(server, rmcResponseBody)
	rmcResponse.ProtocolID = matchmake_extension.ProtocolID
	rmcResponse.MethodID = matchmake_extension.MethodAutoMatchmakeWithParamPostpone
	rmcResponse.CallID = callID

	return rmcResponse, 0
}
