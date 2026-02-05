Add 3 new tools into internal mcp:
1. run build assistant
2. run validator assistant
3. run helper assistant
All of them consume a strict json input:
1. build assistant:
```json
{
    "locationId": "",
    "userIntent": "",
    "mentionedNodes": "",
    "userAcceptanceCriteria": ""
}
```
2. validator assistant:
```json
{
    "userIntent": "",
    "userAcceptanceCriteria": "",
    "automation": {} //automation object
}
```
3. helper assistant:
```json
{
    "userIntent": "",
    "userProblem": "",
    "context":[
        {
            "name": "",
            "value": "",
        }
    ]
}
```