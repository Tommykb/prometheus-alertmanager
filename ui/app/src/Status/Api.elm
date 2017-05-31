module Status.Api exposing (getStatus)

import Utils.Api exposing (baseUrl, send, get)
import Utils.Types exposing (ApiData)
import Status.Types exposing (StatusResponse, VersionInfo, MeshStatus, MeshPeer)
import Json.Decode exposing (Decoder, map2, string, field, at, list, int)


getStatus : (ApiData StatusResponse -> msg) -> Cmd msg
getStatus msg =
    let
        url =
            String.join "/" [ baseUrl, "status" ]

        request =
            get url decodeStatusResponse
    in
        Cmd.map msg <| send request


decodeStatusResponse : Decoder StatusResponse
decodeStatusResponse =
    field "data" decodeData


decodeData : Decoder StatusResponse
decodeData =
    Json.Decode.map4 StatusResponse
        (field "config" string)
        (field "uptime" string)
        (field "versionInfo" decodeVersionInfo)
        (field "meshStatus" decodeMeshStatus)


decodeVersionInfo : Decoder VersionInfo
decodeVersionInfo =
    Json.Decode.map6 VersionInfo
        (field "branch" string)
        (field "buildDate" string)
        (field "buildUser" string)
        (field "goVersion" string)
        (field "revision" string)
        (field "version" string)


decodeMeshStatus : Decoder MeshStatus
decodeMeshStatus =
    Json.Decode.map3 MeshStatus
        (field "name" string)
        (field "nickName" string)
        (field "peers" (list decodeMeshPeer))


decodeMeshPeer : Decoder MeshPeer
decodeMeshPeer =
    Json.Decode.map3 MeshPeer
        (field "name" string)
        (field "nickName" string)
        (field "uid" int)
