# mac-models
Apple Mac, Model name / Model Identifier Table and functions

### *You can find your mac model by code now.*

The mac-models.json file maps model identifiers to model names for all Mac computers. The program can get the model identifier by reading the hw.model information. It can then use this identifier to look up the corresponding model name in the mapping table contained in mac-models.json.

```golang 
        jsonFile := "mac-models.json"
	if !fileExists(jsonFile) {
	  GenerateModelsJson()
	}
	
        modelName := GetMacInfo()
	fmt.Println("Current Mac model: " + modelName)
```

Author: ihugang@gmail.com 2023

*Data come from https://everymac.com*
