module Views.AlertList.FilterBar exposing (view)

import Html exposing (Html, Attribute, div, span, input, text, button, i, small)
import Html.Attributes exposing (value, class, style, disabled, id)
import Html.Events exposing (onClick, onInput, on, keyCode)
import Utils.Filter exposing (Matcher)
import Views.AlertList.Types exposing (AlertListMsg(..))
import Types exposing (Msg(..))
import Json.Decode as Json


viewMatcher : Matcher -> Html Msg
viewMatcher matcher =
    div [ class "col col-auto", style [ ( "padding", "5px" ) ] ]
        [ button
            [ class "btn btn-outline-info"
            , onClick (DeleteFilterMatcher False matcher |> MsgForAlertList)
            ]
            [ text <| Utils.Filter.stringifyMatcher matcher
            , text " ×"
            ]
        ]


lastElem : List a -> Maybe a
lastElem =
    List.foldl (Just >> always) Nothing


viewMatchers : List Matcher -> List (Html Msg)
viewMatchers matchers =
    matchers
        |> List.map viewMatcher


onKeydown : Int -> Msg -> Attribute Msg
onKeydown key msg =
    on "keydown"
        (Json.map
            (\k ->
                if k == key then
                    msg
                else
                    Noop
            )
            keyCode
        )


view : List Matcher -> String -> Html Msg
view matchers matcherText =
    let
        className =
            if matcherText == "" then
                ""
            else
                case maybeMatcher of
                    Just _ ->
                        "has-success"

                    Nothing ->
                        "has-danger"

        maybeMatcher =
            Utils.Filter.parseMatcher matcherText

        onKeydownAttr =
            maybeMatcher
                |> Maybe.map (AddFilterMatcher True >> MsgForAlertList)
                |> Maybe.withDefault Noop
                |> onKeydown 13

        onClickAttr =
            maybeMatcher
                |> Maybe.map (AddFilterMatcher True >> MsgForAlertList)
                |> Maybe.withDefault Noop
                |> onClick

        isDisabled =
            maybeMatcher == Nothing
    in
        div
            [ class "row no-gutters align-items-start" ]
            (viewMatchers matchers
                ++ [ div
                        [ class ("col form-group " ++ className)
                        , style
                            [ ( "padding", "5px" )
                            , ( "min-width", "200px" )
                            , ( "max-width", "500px" )
                            ]
                        ]
                        [ div [ class "input-group" ]
                            [ input
                                [ id "custom-matcher", class "form-control", value matcherText, onKeydownAttr, onInput (UpdateMatcherText >> MsgForAlertList) ]
                                []
                            , span
                                [ class "input-group-btn" ]
                                [ button [ class "btn btn-primary", disabled isDisabled, onClickAttr ] [ text "Add" ] ]
                            ]
                        , small [ class "form-text text-muted" ]
                            [ text "Custom matcher, e.g."
                            , button [ class "btn btn-link btn-sm align-baseline", onClick (UpdateMatcherText "env=\"production\"" |> MsgForAlertList) ] [ text "env=\"production\"" ]
                            ]
                        ]
                   ]
            )
