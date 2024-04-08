package utils

import "github.com/gofiber/fiber/v2/utils"

func IsIPInterface(input interface{}) bool {
	if IsString(input) {
		return IsIP(input.(string))
	}
	return false
}

func IsIP(input string) bool {
	return utils.IsIPv4(input) || utils.IsIPv6(input)
}