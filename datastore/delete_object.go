package datastore

import (
	"github.com/PretendoNetwork/nex-go"
	common_globals "github.com/PretendoNetwork/nex-protocols-common-go/globals"
	datastore "github.com/PretendoNetwork/nex-protocols-go/datastore"
	datastore_types "github.com/PretendoNetwork/nex-protocols-go/datastore/types"
)

func deleteObject(err error, client *nex.Client, callID uint32, param *datastore_types.DataStoreDeleteParam) uint32 {
	if commonDataStoreProtocol.getObjectInfoByDataIDHandler == nil {
		common_globals.Logger.Warning("GetObjectInfoByDataID not defined")
		return nex.Errors.Core.NotImplemented
	}

	if commonDataStoreProtocol.verifyObjectPermissionHandler == nil {
		common_globals.Logger.Warning("VerifyObjectPermission not defined")
		return nex.Errors.Core.NotImplemented
	}

	if commonDataStoreProtocol.deleteObjectByDataIDWithPasswordHandler == nil {
		common_globals.Logger.Warning("DeleteObjectByDataIDWithPassword not defined")
		return nex.Errors.Core.NotImplemented
	}

	if err != nil {
		common_globals.Logger.Error(err.Error())
		return nex.Errors.DataStore.Unknown
	}

	metaInfo, errCode := commonDataStoreProtocol.getObjectInfoByDataIDHandler(param.DataID)
	if errCode != 0 {
		return errCode
	}

	errCode = commonDataStoreProtocol.verifyObjectPermissionHandler(metaInfo.OwnerID, client.PID(), metaInfo.DelPermission)
	if errCode != 0 {
		return errCode
	}

	errCode = commonDataStoreProtocol.deleteObjectByDataIDWithPasswordHandler(param.DataID, param.UpdatePassword)
	if errCode != 0 {
		return errCode
	}

	rmcResponse := nex.NewRMCResponse(datastore.ProtocolID, callID)
	rmcResponse.SetSuccess(datastore.MethodDeleteObject, nil)

	rmcResponseBytes := rmcResponse.Bytes()

	responsePacket, _ := nex.NewPacketV1(client, nil)

	responsePacket.SetVersion(1)
	responsePacket.SetSource(0xA1)
	responsePacket.SetDestination(0xAF)
	responsePacket.SetType(nex.DataPacket)
	responsePacket.SetPayload(rmcResponseBytes)

	responsePacket.AddFlag(nex.FlagNeedsAck)
	responsePacket.AddFlag(nex.FlagReliable)

	commonDataStoreProtocol.server.Send(responsePacket)

	return 0
}
