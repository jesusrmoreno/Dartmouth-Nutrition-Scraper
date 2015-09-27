package constants

// AvailableSIDSRequest is the Request String to be used in getting the
// Available sids
const AvailableSIDSRequest string = `
  {
    "service": "",
    "method": "get_available_sids",
    "id": 1,
    "params": [
      null,
      "{ \"remoteProcedure\":\"get_available_sids\" }"
    ]
  }
`

// GetSIDSRequest ...
const GetSIDSRequest string = `{
  "service": "",
  "method": "create_context",
  "id": 2,
  "params": ["%s"]
}`

// GetMenuListRequest ...
const GetMenuListRequest string = `{
  "service": "",
  "method": "get_webmenu_list",
  "id": 5,
  "params": [{
    "sid":"%s"
  }, "{\"remoteProcedure\":\"get_webmenu_list\"}"]
}`
