module github.com/johnkerl/goffl

go 1.25

// When pgpg is not published, use local replaces (e.g. sibling repo):
replace github.com/johnkerl/pgpg/lib => ../pgpg/lib

replace github.com/johnkerl/pgpg/generators/go => ../pgpg/generators/go

require github.com/johnkerl/pgpg/lib v0.0.0

require github.com/johnkerl/pgpg/generators/go v0.0.0-20260222002111-906fc64c96a8 // indirect
