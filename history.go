package gocache

// store requests in file
// with following fmt:
// {KeyRequest uint8}+{operation data: writeOneRequest | eraseOneRequest etc.}

// this can be useful for logging and reconstructing state from past history
