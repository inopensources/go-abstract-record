package models

import (
	"database/sql"
	"fmt"
	"reflect"
)

type AbstractModel struct {
	Object interface{}
	db *sql.DB
}

func (a AbstractModel) All(table string, pagination Pagination) (*sql.Rows, error) {
	rows, err := a.db.Query(mountSql(table, pagination))
	a.db.Close()
	return rows, err
}

func (a AbstractModel) FindByKey(table string, condictions []string, keys []string, operations []string, values []string) (*sql.Rows, error) {
	rows, err := a.db.Query(mountQuery(table, condictions, keys, operations, values))
	a.db.Close()
	return rows, err
}

func (a AbstractModel) Count(table string) (int, error) {
	rows, err := a.db.Query("SELECT COUNT(*) FROM " + table)
	if err != nil {
		return 0, err
	}
	a.db.Close()

	var count int
	for rows.Next() {
		err = rows.Scan(&count)
		if err != nil {
			fmt.Printf(err.Error())
		}
	}

	return count + 1, nil
}

func (a AbstractModel) Create(table string) (sql.Result, error) {
	insert := fmt.Sprint("INSERT INTO ", table, " (")

	elements := reflect.ValueOf(a.Object).Elem()
	typeOfT := elements.Type()
	t := reflect.TypeOf(a.Object)

	for i := 0; i < elements.NumField(); i++ {
		ff, _ := t.Elem().FieldByName(typeOfT.Field(i).Name)

		if (i + 1) == elements.NumField() { // ultimo campo
			insert = fmt.Sprint(insert, ff.Tag, ") VALUES (")
		} else {
			insert = fmt.Sprint(insert, ff.Tag, ", ")
		}
	}

	for i := 0; i < elements.NumField(); i++ {
		f := elements.Field(i)

		switch f.Type().String() {
		case "string":
			insert = fmt.Sprint(insert, "'", f.Interface(), "'", checaSeEUltimaInstrucao(i, elements.NumField(), ")"))
		case "bool":
			insert = fmt.Sprint(insert, "'", f.Interface(), "'", checaSeEUltimaInstrucao(i, elements.NumField(), ")"))
		case "int":
			insert = fmt.Sprint(insert, f.Interface(), checaSeEUltimaInstrucao(i, elements.NumField(), ")"))
		case "float":
			insert = fmt.Sprint(insert, f.Interface(), checaSeEUltimaInstrucao(i, elements.NumField(), ")"))
		default:
			insert = fmt.Sprint(insert, "'", f.Interface(), "'", checaSeEUltimaInstrucao(i, elements.NumField(), ")"))
		}
	}

	result, err := a.db.Exec(insert)
	a.db.Close()
	return result, err
}

func (a AbstractModel) Update(table string, keys []string) (sql.Result, error) {
	update := fmt.Sprint("UPDATE ", table, " SET ")

	elements := reflect.ValueOf(a.Object).Elem()
	typeOfT := elements.Type()
	t := reflect.TypeOf(a.Object)

	for i := 0; i < elements.NumField(); i++ {
		ff, _ := t.Elem().FieldByName(typeOfT.Field(i).Name)
		f := elements.Field(i)

		contains, _ := Contains(keys, string(ff.Tag))
		if !contains { // Só pode entrar aqui quem NAO for chave
			update = checkTypes(f, ff, i, elements, update, " = ", ",")
		}
	}

	naoEntrouAinda := true
	for i := 0; i < elements.NumField(); i++ {
		ff, _ := t.Elem().FieldByName(typeOfT.Field(i).Name)
		f := elements.Field(i)

		contains, _ := Contains(keys, string(ff.Tag))
		if contains { // Só pode entrar aqui quem FOR chave
			if naoEntrouAinda {
				naoEntrouAinda = false // seta para falso para não entrar duas vezes no WHERE
				update = checkTypes(f, ff, i, elements, update+" WHERE ", " = ", "")
			} else {
				update = checkTypes(f, ff, i, elements, update+" AND ", " = ", "")
			}
		}
	}

	result, err := a.db.Exec(update)
	a.db.Close()
	return result, err
}

func (a AbstractModel) Delete(table string, keys []string, values []string) (sql.Result, error) {
	delete := fmt.Sprint("DELETE FROM ", table)

	elements := reflect.ValueOf(a.Object).Elem()
	typeOfT := elements.Type()
	t := reflect.TypeOf(a.Object)

	naoEntrouAinda := true
	for i := 0; i < elements.NumField(); i++ {
		ff, _ := t.Elem().FieldByName(typeOfT.Field(i).Name)
		f := elements.Field(i)

		contains, index := Contains(keys, string(ff.Tag))
		if contains { // Só pode entrar aqui quem FOR chave
			if naoEntrouAinda {
				naoEntrouAinda = false // seta para falso para não entrar duas vezes no WHERE
				delete = fmt.Sprint(delete, " WHERE ")
			} else {
				delete = fmt.Sprint(delete, " AND ")
			}

			switch f.Type().String() {
			case "string":
				delete = fmt.Sprint(delete, ff.Tag, " = ", "'", f.Interface(), "'", values[index])
			case "bool":
				delete = fmt.Sprint(delete, ff.Tag, " = ", "'", f.Interface(), "'", values[index])
			case "int":
				delete = fmt.Sprint(delete, ff.Tag, " = ", f.Interface(), values[index])
			case "float":
				delete = fmt.Sprint(delete, ff.Tag, " = ", f.Interface(), values[index])
			default:
				delete = fmt.Sprint(delete, ff.Tag, " = ", "'", f.Interface(), "'", values[index])
			}
		}
	}

	result, err := a.db.Exec(delete)
	a.db.Close()
	return result, err
}

func mountSql(table string, pagination Pagination) string {
	sql := fmt.Sprint(" SELECT * FROM ", table)
	if pagination.Where != "" {
		sql = fmt.Sprint(sql, " ", pagination.Where)
	}

	if pagination.GroupBy != "" {
		sql = fmt.Sprint(sql, " ", pagination.GroupBy)
	}

	if pagination.OrderBy != "" {
		sql = fmt.Sprint(sql, " ", pagination.OrderBy)
	} else {
		sql = fmt.Sprint(sql, " ORDER BY 1 ")
	}

	sql = fmt.Sprint(sql, " OFFSET ", ((pagination.MorePerPage * pagination.Page) - pagination.MorePerPage), " ROWS FETCH NEXT ", pagination.MorePerPage, " ROWS ONLY ")

	return sql
}

/**
		condictions tem que ser sempre um vetor menor que os demais, pois a primeira posicao sera WHERE
 */
func mountQuery(table string, condictions []string, keys []string, operations []string, values []string) string {
	sql := fmt.Sprint(" SELECT * FROM ", table)
	for i := 0; i < len(keys); i++ {
		if i == 0 {
			sql = fmt.Sprint(sql, " WHERE ", keys[i], " ", operations[i], " ", values[i])
		} else {
			sql = fmt.Sprint(sql, condictions[i-1], " ", keys[i], " ", operations[i], " ", values[i])
		}
	}
	return sql
}

func checaSeEUltimaInstrucao(index int, total int, final string) string {
	if (index + 1) == total {
		return final
	}
	return ", "
}

/**
	Contains em um ARRAY
 */
func Contains(s []string, e string) (bool, int) {
	for i, a := range s {
		if a == e {
			return true, i
		}
	}
	return false, 0
}

func checkTypes(f reflect.Value, ff reflect.StructField, i int, elements reflect.Value, sql string, operation string, final string) string {
	switch f.Type().String() {
	case "string":
		sql = fmt.Sprint(sql, ff.Tag, operation, "'", f.Interface(), "'", checaSeEUltimaInstrucao(i, elements.NumField(), final))
	case "bool":
		sql = fmt.Sprint(sql, ff.Tag, operation, "'", f.Interface(), "'", checaSeEUltimaInstrucao(i, elements.NumField(), final))
	case "int":
		sql = fmt.Sprint(sql, ff.Tag, operation, f.Interface(), checaSeEUltimaInstrucao(i, elements.NumField(), final))
	case "float":
		sql = fmt.Sprint(sql, ff.Tag, operation, f.Interface(), checaSeEUltimaInstrucao(i, elements.NumField(), final))
	default:
		sql = fmt.Sprint(sql, ff.Tag, operation, "'", f.Interface(), "'", checaSeEUltimaInstrucao(i, elements.NumField(), final))
	}
	return sql
}
