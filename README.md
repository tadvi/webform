# webform
Go (golang) package webform fills a struct with form values

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

Check doc.go for details.

Copyright (C) Tad Vizbaras, MIT license.
