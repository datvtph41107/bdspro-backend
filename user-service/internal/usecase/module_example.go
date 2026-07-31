package usecase

import (
	"fmt"
)

// ExampleUsage ví dụ cách sử dụng ModuleService
func ExampleUsage() {
	// Ví dụ sử dụng ModuleService sau khi đã load modules

	// Giả sử có instance của ModuleService
	// moduleService := app.ModuleService

	// Lấy module theo ID
	// module := moduleService.GetModuleByID(1)
	// if module != nil {
	//     fmt.Printf("Module found by ID: %s\n", module.Name)
	// }

	// Lấy module theo Key
	// module = moduleService.GetModuleByKey("USER_MANAGEMENT")
	// if module != nil {
	//     fmt.Printf("Module found by Key: %s\n", module.Name)
	// }

	// Lấy module theo Code
	// module = moduleService.GetModuleByCode("UM001")
	// if module != nil {
	//     fmt.Printf("Module found by Code: %s (ID: %d)\n", module.Name, module.ID)
	// }

	// Lấy tất cả modules theo ID
	// allModules := moduleService.GetAllModules()
	// fmt.Printf("Total modules loaded: %d\n", len(allModules))

	// Lấy tất cả modules theo Key
	// allModulesByKey := moduleService.GetAllModulesByKey()
	// fmt.Printf("Total modules by key: %d\n", len(allModulesByKey))

	// Lấy tất cả modules theo Code
	// allModulesByCode := moduleService.GetAllModulesByCode()
	// fmt.Printf("Total modules by code: %d\n", len(allModulesByCode))

	fmt.Println("ModuleService example usage - check comments above")
}
