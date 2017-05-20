module Views.Shared.SilencePreview exposing (view)

import Silences.Types exposing (Silence)
import Html exposing (Html, div, text)
import Html.Attributes exposing (class)
import Utils.Types exposing (ApiResponse(Initial, Success, Loading, Failure))
import Views.Shared.AlertListCompact
import Utils.Views exposing (error, loading)


view : Silence -> Html msg
view { silencedAlerts } =
    case silencedAlerts of
        Success alerts ->
            if List.isEmpty alerts then
                div [ class "w-100 mt-3" ] [ text "No matches" ]
            else
                div [ class "w-100 mt-3" ] [ Views.Shared.AlertListCompact.view alerts ]

        Initial ->
            text ""

        Loading ->
            loading

        Failure e ->
            error e
