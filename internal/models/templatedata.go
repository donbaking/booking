package models

//建立一個TemplateData的 struct用來接收前端給的資料或是把資料給前端
type TemplateData struct{
	// StringMap is a map with string keys and string values
	StringMap map[string]string
	// IntMap is a map with string keys and int values
	IntMap map[string]int
	// FloatMap is a map with string keys and float32 values
	FloatMap map[string]float32
	//Data map在不確定values的狀態下可以使用interfaceP{}去存取不定的values datatype
	Data map[string]interface{}
	//token 
	CSRFToken string
	//success message and error message
	Flash string
	Warning string
	Error string
}