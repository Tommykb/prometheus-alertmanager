module Alerts.Api exposing (..)

import Alerts.Types exposing (Alert, RouteOpts, Block, AlertGroup)
import Json.Decode as Json exposing (..)
import Utils.Api exposing (iso8601Time)
import Utils.Types exposing (ApiData)
import Utils.Filter exposing (Filter, generateQueryString)


fetchAlerts : String -> Filter -> Cmd (ApiData (List Alert))
fetchAlerts baseUrl filter =
    let
        url =
            String.join "/" [ baseUrl, "alerts" ++ (generateQueryString filter) ]
    in
        Utils.Api.send (Utils.Api.get url alertsDecoder)


alertsDecoder : Json.Decoder (List Alert)
alertsDecoder =
    Json.list alertDecoder
        -- populate alerts with ids:
        |> Json.map (List.indexedMap (toString >> (|>)))
        |> field "data"


{-| TODO: decode alert id when provided
-}
alertDecoder : Json.Decoder (String -> Alert)
alertDecoder =
    Json.map6 Alert
        (Json.maybe (field "annotations" (Json.keyValuePairs Json.string))
            |> andThen (Maybe.withDefault [] >> Json.succeed)
        )
        (field "labels" (Json.keyValuePairs Json.string))
        (Json.maybe (Json.at [ "status", "silencedBy", "0" ] Json.string))
        (Json.maybe (Json.at [ "status", "inhibitedBy", "0" ] Json.string)
            |> Json.map ((/=) Nothing)
        )
        (field "startsAt" iso8601Time)
        (field "generatorURL" Json.string)
