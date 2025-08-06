BEGIN {
    iota = 0  # Current token number
    watch = 0 # Be sensitive to lines?
}

# Begin watching only within const () blocks
/^const/ {watch = 1}

/\tItem/ {
    if (watch == 1) {
        # Trim the Item prefix
        sub(/^Item/, "", $1)

        # Generate the output we want to feed into YACC
        printf("%%token <num> %s\t%d\n", toupper($1), iota)

        # Increment the token number by one (i.e. just as iota does in Go)
        iota++
    }
}

# Once we leave the block, stop watching!
/^\)/ {
    if (watch == 1) {
        watch = false
    }
}
