package repository

import (
	"fmt"
	"strconv"
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
			%s
		FROM "%ss"
		`, queryFieldsBuilder(false, fields), entity))
}

func getOneQuery(entity string, fields ...string) string {
	return removeEscapeChar(
		fmt.Sprintf(`
		SELECT 
			%s
		FROM "%ss"
		WHERE "%sID"=$1
	`, queryFieldsBuilder(false, fields), entity, entity))
}

func createOneQuery(entity string, fields ...string) string {
	return removeEscapeChar(
		fmt.Sprintf(`
		INSERT INTO "%ss"
			(%s)
		VALUES (%s)
		RETURNING "%sID"`, entity, queryFieldsBuilder(false, fields), queryValuesPlaceholder(len(fields), 1), entity))
}

func AddManyToMany(firstEntity string, secondEntity string, extraFields ...string) string {
	return removeEscapeChar(
		fmt.Sprintf(`
		INSERT INTO "%ss_%ss"
		("%sID", "%sID"%s)
		VALUES (%s)
		`, firstEntity, secondEntity, firstEntity, secondEntity, queryFieldsBuilder(true, extraFields), queryValuesPlaceholder(len(extraFields)+2, 1)))
}

func AddOrUpdateManyToMany(firstEntity string, secondEntity string, extraFields ...string) string {
	return removeEscapeChar(
		fmt.Sprintf(`
		INSERT INTO "%ss_%ss"
			("%sID", "%sID"%s)
		VALUES (%s)
		ON CONFLICT("%sID", "%sID")
		DO UPDATE SET %s
		`, firstEntity, secondEntity, firstEntity, secondEntity, queryFieldsBuilder(true, extraFields), queryValuesPlaceholder(len(extraFields)+2, 1),
			firstEntity, secondEntity, querySetFieldsBuilder(extraFields, 3))) //update parameters

}

func removeEscapeChar(str string) string {
	return strings.Replace(str, `\`, "", -1)
}

func queryValuesPlaceholder(valuesCount int, startFrom int) string {
	var b []string
	for i := startFrom; i < startFrom+valuesCount; i++ {
		b = append(b, fmt.Sprintf("$%d", i))
	}
	return strings.Join(b, ", ")
}

func queryFieldsBuilder(startWithComma bool, fields []string) string {
	var b strings.Builder
	if startWithComma {
		b.WriteString(", ")
	}
	commaSeparatedFields := strings.Join(fields, `", "`)
	commaSeparatedFields = strconv.Quote(commaSeparatedFields)
	fmt.Fprintf(&b, commaSeparatedFields)
	return b.String()
}

func querySetFieldsBuilder(fields []string, startFrom int) string {
	fieldsCopy := append([]string(nil), fields...)
	for i := 0; i < len(fieldsCopy); i++ {
		fieldsCopy[i] = strconv.Quote(fieldsCopy[i]) + fmt.Sprintf("=$%d", startFrom+i)
	}
	return strings.Join(fieldsCopy, ",")
}
