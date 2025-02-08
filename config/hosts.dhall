let Entries
    : Type
    = { IP : Text, Aliases : List Text }

let Group
    : Type
    = { Name : Text, Entries : List Entries }

let Config
    : Type
    = { Groups : List Group }

let config
    : Config
    = { Groups =
        [ { Name = "sample"
          , Entries =
            [ { IP = "127.0.0.1", Aliases = [ "localhost" ] }
            , { IP = "255.255.255.255", Aliases = [ "broadcasthost" ] }
            , { IP = "::1", Aliases = [ "localhost" ] }
            ]
          }
        ]
      }

in  config
