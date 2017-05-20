module Utils.Views exposing (..)

import Html exposing (..)
import Html.Attributes exposing (..)
import Html.Events exposing (onCheck, onInput, onClick)
import Http exposing (Error(..))


tab : tab -> tab -> (tab -> msg) -> List (Html msg) -> Html msg
tab tab currentTab msg content =
    li [ class "nav-item" ]
        [ if tab == currentTab then
            span [ class "nav-link active" ] content
          else
            a [ class "nav-link", onClick (msg tab) ] content
        ]


labelButton : Maybe msg -> String -> Html msg
labelButton maybeMsg labelText =
    case maybeMsg of
        Nothing ->
            span
                [ class "btn btn-sm bg-faded btn-secondary mr-2 mb-2" ]
                [ text labelText ]

        Just msg ->
            button
                [ class "btn btn-sm bg-faded btn-secondary mr-2 mb-2"
                , onClick msg
                ]
                [ span [ class "text-muted" ] [ text labelText ] ]


iconButtonMsg : String -> String -> msg -> Html msg
iconButtonMsg classString icon msg =
    a [ class classString, onClick msg ]
        [ i [ class <| "fa fa-3 " ++ icon ] []
        ]


checkbox : String -> Bool -> (Bool -> msg) -> Html msg
checkbox name status msg =
    label [ class "f6 dib mb2 mr2 d-flex align-items-center" ]
        [ input [ type_ "checkbox", checked status, onCheck msg ] []
        , span [ class "pl-2" ] [ text <| " " ++ name ]
        ]


validatedFormField : String -> Result ( String, String ) String -> String -> (String -> msg) -> Html msg
validatedFormField labelText validatedString classes msg =
    case validatedString of
        Ok inputValue ->
            div [ class <| "d-flex flex-column form-group has-success " ++ classes ]
                [ label [] [ strong [] [ text labelText ] ]
                , input [ value inputValue, onInput msg, class "form-control form-control-success" ] []
                ]

        Err ( inputValue, error ) ->
            div [ class <| "d-flex flex-column form-group has-danger " ++ classes ]
                [ label [] [ strong [] [ text labelText ] ]
                , input [ value inputValue, onInput msg, class "form-control form-control-danger" ] []
                , div [ class "form-control-feedback" ] [ text error ]
                ]


formField : String -> String -> String -> (String -> msg) -> Html msg
formField labelText content classes msg =
    div [ class <| "d-flex flex-column " ++ classes ]
        [ label [] [ strong [] [ text labelText ] ]
        , input [ value content, onInput msg ] []
        ]


validatedTextField : String -> Result ( String, String ) String -> String -> (String -> msg) -> Html msg
validatedTextField labelText validatedString classes msg =
    case validatedString of
        Ok inputValue ->
            div [ class <| "d-flex flex-column form-group has-success " ++ classes ]
                [ label [] [ strong [] [ text labelText ] ]
                , textarea [ value inputValue, onInput msg, class "form-control form-control-success" ] []
                ]

        Err ( inputValue, error ) ->
            div [ class <| "d-flex flex-column form-group has-danger " ++ classes ]
                [ label [] [ strong [] [ text labelText ] ]
                , textarea [ value inputValue, onInput msg, class "form-control form-control-danger" ] []
                , div [ class "form-control-feedback" ] [ text error ]
                ]


textField : String -> String -> String -> (String -> msg) -> Html msg
textField labelText content classes msg =
    div [ class <| "d-flex flex-column " ++ classes ]
        [ label [] [ strong [] [ text labelText ] ]
        , textarea [ value content, onInput msg ] []
        ]


buttonLink : String -> String -> String -> msg -> Html msg
buttonLink icon link color msg =
    a [ class <| "" ++ color, href link, onClick msg ]
        [ i [ class <| "" ++ icon ] []
        ]


formInput : String -> String -> (String -> msg) -> Html msg
formInput inputValue classes msg =
    Html.input [ class <| "w-100 " ++ classes, value inputValue, onInput msg ] []


loading : Html msg
loading =
    div []
        [ i [ class "fa fa-cog fa-spin fa-3x fa-fw" ] []
        , span [ class "sr-only" ] [ text "Loading..." ]
        ]


error : Http.Error -> Html msg
error err =
    let
        msg =
            case err of
                Timeout ->
                    "timeout exceeded"

                NetworkError ->
                    "network error"

                BadStatus resp ->
                    resp.status.message ++ " " ++ resp.body

                BadPayload err resp ->
                    -- OK status, unexpected payload
                    "unexpected response from api" ++ err

                BadUrl url ->
                    "malformed url: " ++ url
    in
        div []
            [ p [] [ text <| "Error: " ++ msg ]
            ]
