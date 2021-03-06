package main

import (
	"fmt"
	"reflect"
)

type DeviceStatus uint8

const (
	Offline DeviceStatus = 0
	Online DeviceStatus = 1
)

type Category struct {
	ID int32
	Name string `json:"name" value:"空调"`
	Description  string
}

func main() {
	var category Category
	var i int32 = -5
	

	// 类型（Type）指的是系统原生数据类型，如 int、string、bool、float32 等类型，以及使用 type 关键字定义的类型
	// 种类（Kind）指的是对象归属的品种. 在go语言中是固定定义好的，可以通过查看Kind定义知道全部Kind
	typeInfo, kindInfo := reflectTypeAndKind(category)
	fmt.Printf("struct type reflect info.            name:%s, kind:%s \n", typeInfo, kindInfo)
	enumerateFieldByReflect(category)

	typeOfCagegory := reflect.TypeOf(category)
	typeInfo, kindInfo = reflectTypeAndKind(typeOfCagegory)
	fmt.Printf("type reflect info of reflect.        name:%s, kind:%s \n", typeInfo, kindInfo)
	enumerateFieldByReflect(category)
	fmt.Println("categroy is not initialized.")
	getAndUpdateFieldValueByReflect(&category)
	category = Category{ID:001, Name: "空调", Description: "空调大类"}
	fmt.Printf("---updateFieldValue by reflect.   rawValue:%v \n\n", category)
	fmt.Println("categroy is initialized.")
	fmt.Printf("---enumerateFieldValue by reflect.   rawValue:%v \n", category)
	enumerateFieldByReflect(category)
	fmt.Printf("---start of updateFieldValue by reflect.   rawValue:%v \n", category)
	getAndUpdateFieldValueByReflect(&category)
	fmt.Printf("---end of updateFieldValue by reflace.   newValue:%v \n\n", category)

	typeInfo, kindInfo = reflectTypeAndKind(i)
	fmt.Printf("primitive type reflect info.         name:%s, kind:%v \n", typeInfo, kindInfo)
	enumerateFieldByReflect(i)
	getAndUpdateFieldValueByReflect(i)

	myFunc := reflectTypeAndKind
	typeInfo, kindInfo = reflectTypeAndKind(myFunc)
	fmt.Printf("func type reflect info.              name:%s, kind:%v \n", typeInfo, kindInfo)
	enumerateFieldByReflect(myFunc)
	getAndUpdateFieldValueByReflect(myFunc)

	var devStatus DeviceStatus = Online
	typeInfo, kindInfo = reflectTypeAndKind(devStatus)
	fmt.Printf("enum type reflect info.              name:%s, kind:%v \n", typeInfo, kindInfo)
	enumerateFieldByReflect(devStatus)
	getAndUpdateFieldValueByReflect(devStatus)
}

func reflectTypeAndKind(x interface{}) (typeInfo string, kindInfo reflect.Kind) {
	reflectObject := reflect.TypeOf(x)
	return reflectObject.Name(), reflectObject.Kind()
}

func enumerateFieldByReflect(x interface{}) {
	reflectObject := reflect.TypeOf(x)
	reflectValue := reflect.ValueOf(x)

	// NumField returns a struct type's field count.
	// It panics if the type's Kind is not Struct.
	if reflectObject.Kind() == reflect.Struct {
		num := reflectObject.NumField();
		fmt.Printf("there are %d fields\n", num)
		for i := 0; i < num; i++ {
			// 获取每个属性的结构体字段类型
			fieldType := reflectObject.Field(i)
			
			filedValue := reflectValue.Field(i)
			fieldKind := filedValue.Kind()

			var rawIntValue int32
			var rawStrValue string
			switch fieldKind {
			    case reflect.Int32:
				//获取64位的值，强制类型转换为int类型
				    rawIntValue = filedValue.Interface().(int32)
			    case reflect.String:
				    rawStrValue = filedValue.Interface().(string)
			    default:
		    } 
		 	// 输出属性名和tag
			fmt.Printf("%dth, name: %v,  type %v, tag: '%v', currentIntValue:%v, currentStrValue:%v \n", 
				i, fieldType.Name, fieldType.Type, fieldType.Tag, rawIntValue, rawStrValue)
		}
	  
		// 通过字段名, 找到字段类型信息
		if nameField, ok := reflectObject.FieldByName("Name"); ok {
			// 从tag中取出需要的tag
			fmt.Println(nameField.Tag.Get("json"), nameField.Tag.Get("value"))
		} else {
			fmt.Println("no name filed")
		}
		// It panics if v's Kind is not struct.
		fmt.Println("不存在的结构体成员:", reflect.ValueOf(x).FieldByName("").IsValid())
	} else {
		fmt.Println("non-struct, no filed")
	}
}
/*
* 更新值，必须传入指针类型， 反射只能修改可以导出的属性（也就是大写字母开始的属性）
*/
func getAndUpdateFieldValueByReflect(x interface{}) {
	reflectObject := reflect.TypeOf(x)
	reflectValue := reflect.ValueOf(x)
	reflectKind := reflectObject.Kind()
	actualReflectKind := reflect.Ptr
	
	if reflectKind == reflect.Ptr {
		actualReflectKind =  reflect.Struct
		fmt.Printf("ptr type: %T\n", x)
		reflectObject = reflectObject.Elem()
		actualReflectKind = reflectObject.Kind()
		reflectValue = reflectValue.Elem()

	} 
	// NumField returns a struct type's field count.
	// It panics if the type's Kind is not Struct.
	if actualReflectKind == reflect.Struct {
		num := reflectObject.NumField();
		fmt.Printf("there are %d fields, reflectKind:%v, actualReflectKind:%v\n",
		    num, reflectKind, actualReflectKind)
		for i := 0; i < num; i++ {
			// 获取每个属性的结构体字段类型
			fieldType := reflectObject.Field(i)
			filedValue := reflectValue.Field(i)
			fieldKind := filedValue.Kind()
			// 输出属性名和tag
			fmt.Printf("%dth, name: %v,  type %v, tag: '%v', filedValue:%v, fieldValueInterface:%v \n",
			    i, fieldType.Name, fieldType.Type, fieldType.Tag, filedValue, filedValue.Interface())

			switch fieldKind {
				case reflect.Int32:					
					newIntValue := 999 + filedValue.Interface().(int32)
					// 强制将int32转换为int64，否则无法传入参数
				    filedValue.SetInt(int64(newIntValue))
				case reflect.String:
					newStrValue := "aaa-" + filedValue.Interface().(string)
				    filedValue.SetString(newStrValue)
				default:
			 } 
		}
	  
		// 通过字段名, 找到字段类型信息
		if nameField, ok := reflectObject.FieldByName("Name"); ok {
			fmt.Println(nameField.Tag.Get("json"), nameField.Tag.Get("value"), reflectValue.FieldByName("Name"))
		} else {
			fmt.Println("no name filed")
		}
	} else {
		
		fmt.Printf("non-struct, no filed. reflectKind:%v, rawValue: %v \n", reflectKind, reflectValue)
	}
	
}

