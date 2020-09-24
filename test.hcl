project = "test" 

properties = {
    test = "asdf",
    tata = "asasDF"
}

target "clean" {
    // step "delete" {
    //     dir = "test"
    // }
    step "mkdir" {
        dir = "/tmp/plugin"
    }
}

// target "build" {
//     step "mkdir" {
//         dir = "asd"
//     }
// }
