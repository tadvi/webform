/*
Package webform fills a struct with form values.
It borrows main idea from gorilla/schema which is more complete
solution for complex forms.

webform is simplified idea of the gorilla/schema.
No metadata cache or dotted notation for embedded structs.
Simply take posted form data and map it to structs. That's it.

Only uses url.Values first value of the field. In some cases
web forms are designed to submit fields with multiple values, in that
case this package will not work.

Sets fields on root struct but also on all nested structs.
If finds pointer to nested struct, then creates struct.

---- If two fields of the structs have the same name and
field comes as form field then:

It panics if webform.PanicOnDuplicateName == true.

OR

Only Root field gets set. If two fields are in
different structs embedded in root, then struct that is
embedded first in the root gets field set.

Use struct tag "webform:" to disambiguate the name.

----

Usage is simple. Given this struct:

	type Person struct {
		Name  string
		Phone string
	}

...we can fill it passing a map to the Load() function:

	values := map[string][]string{
		"Name":  {"John"},
		"Phone": {"999-999-999"},
	}
	person := new(Person)
	webform.Decode(person, values)

Typically it will come from a http.Request object and will be of
type url.Values: http.Request.Form or http.Request.MultipartForm:

	func MyHandler(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()

		if err != nil {
			// Handle error
		}

		// r.PostForm is a map of our POST form values
		err := webform.Decode(person, r.PostForm)
		if err != nil {
			// Handle error
		}

		// Do something with person.Name or person.Phone
	}

*/
package webform
