package repository

import (
	"fmt"
	"strings"
)

func checkIfExistsQuery(entity string) string {
	return fmt.Sprintf(`
	SELECT EXISTS(
		SELECT 
			1
		FROM "%ss"
		WHERE "%sID"=$1)
		`, entity, entity)
}

func getAllQuery(entity string, fields ...string) string {
	return removeEscapeChar(
		fmt.Sprintf(`
		SELECT
			%q
		FROM "%ss"
		`, strings.Join(fields, `", "`), entity))
}

func getOneQuery(entity string, fields ...string) string {
	return removeEscapeChar(
		fmt.Sprintf(`
		SELECT 
			%q
		FROM "%ss"
		WHERE "%sID"=$1
	`, strings.Join(fields, `", "`), entity, entity))
}

func createOneQuery(entity string, fields ...string) string {
	return removeEscapeChar(
		fmt.Sprintf(`
		INSERT INTO "%ss"
			(%q)
		VALUES ($1, $2, $3)
		RETURNING "%sID"`, entity, strings.Join(fields, `", "`), entity))
}

func removeEscapeChar(str string) string {
	return strings.Replace(str, `\`, "", -1)
}
