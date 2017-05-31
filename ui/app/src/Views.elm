module Views exposing (..)

import Html exposing (Html, text, div)
import Html.Attributes exposing (class)
import Types exposing (Msg(MsgForSilenceForm, MsgForSilenceView), Model, Route(..))
import Utils.Views exposing (error, loading)
import Views.SilenceList.Views as SilenceList
import Views.SilenceForm.Views as SilenceForm
import Views.AlertList.Views as AlertList
import Views.SilenceView.Views as SilenceView
import Views.NotFound.Views as NotFound
import Views.Status.Views as Status
import Views.NavBar.Views exposing (navBar)


view : Model -> Html Msg
view model =
    div []
        [ navBar model.route
        , div [ class "container pb-4" ]
            [ currentView model ]
        ]


currentView : Model -> Html Msg
currentView model =
    case model.route of
        StatusRoute ->
            Status.view model.status

        SilenceViewRoute silenceId ->
            SilenceView.view model.silenceView |> Html.map MsgForSilenceView

        AlertsRoute filter ->
            AlertList.view model.alertList filter

        SilenceListRoute _ ->
            SilenceList.view model.silenceList

        SilenceFormNewRoute keep ->
            SilenceForm.view Nothing model.silenceForm |> Html.map MsgForSilenceForm

        SilenceFormEditRoute silenceId ->
            SilenceForm.view (Just silenceId) model.silenceForm |> Html.map MsgForSilenceForm

        TopLevelRoute ->
            Utils.Views.loading

        NotFoundRoute ->
            NotFound.view
