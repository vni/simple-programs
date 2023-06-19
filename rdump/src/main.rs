struct Options {
    width: usize,
    show_filename: bool,
    show_addr: bool,
    show_ascii: bool,
}

impl Options {
    fn default() -> Options {
        Options {
            width: 16,
            show_filename: false,
            show_addr: true,
            show_ascii: true,
        }
    }
}

fn dump_line(o: &Options, addr: usize, slice: &[u8]) {
    if o.show_addr {
        print!("{:08x}: ", addr);
    }

    for v in slice {
        print!("{:02x} ", v);
    }

    if o.show_ascii {
        for v in slice {
            if v.is_ascii_graphic() {
                print!("{}", *v as char);
            } else {
                print!(".");
            }
        }
    }

    println!();
}

fn dump_last_line(o: &Options, addr: usize, sz: usize, slice: &[u8]) {
    if o.show_addr {
        print!("{:08x}: ", addr);
    }

    for i in 0..sz {
        if i < slice.len() {
            print!("{:02x} ", slice[i]);
        } else {
            print!("   ");
        }
    }

    if o.show_ascii {
        for i in 0..sz {
            if i < slice.len() {
                if slice[i].is_ascii_graphic() {
                    print!("{}", slice[i] as char);
                } else {
                    print!(".");
                }
            } else {
                print!(" ");
            }
        }
    }

    println!();
}

fn hexdump(o: &Options, v: &[u8]) {
    let sz = v.len();
    let mut i = 0;
    while i < sz - o.width {
        dump_line(&o, i, &v[i..i + o.width]);
        i += o.width;
    }

    if i < sz {
        dump_last_line(&o, i, o.width, &v[i..]);
    }
}

fn main() {
    let mut o = Options::default();

    let args = std::env::args().collect::<Vec<String>>();
    let mut i = 1;
    while i < args.len() {
        if args[i] == "--no-addr" {
            o.show_addr = false;
            i += 1;
            continue;
        } else if args[i] == "--no-ascii" {
            o.show_ascii = false;
            i += 1;
            continue;
        } else if args[i] == "--no-filename" {
            o.show_filename = false;
            i += 1;
            continue;
        } else if args[i] == "--show-addr" {
            o.show_addr = true;
            i += 1;
            continue;
        } else if args[i] == "--show-ascii" {
            o.show_ascii = true;
            i += 1;
            continue;
        } else if args[i] == "--show-filename" {
            o.show_filename = true;
            i += 1;
            continue;
        } else if args[i] == "--width" {
            if i < args.len() - 1 {
                o.width = match args[i + 1].parse::<usize>() {
                    Ok(width) => width,
                    Err(e) => panic!("--width argument isn't number: err: {}", e),
                };
                i += 1;
                continue;
            } else {
                panic!("--width expected a number argument");
            }
        }

        let filename = &args[i];
        if o.show_filename {
            println!("{}", filename);
        }
        let v = std::fs::read(filename);
        if v.is_err() {
            println!("Failed to open {}", filename);
            continue;
        }

        let v = v.unwrap();
        hexdump(&o, &v);

        i += 1;
    }
}
