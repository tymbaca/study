let world = "Hello from a module";;

let rec length full_list = 
    match full_list with 
    | [] -> 0
    | _ :: rem_list -> 1 + length rem_list;;

let () = print_endline world;;


let () = Printf.printf "%d\n" (length [1, 3, 6, 2, 7]);;
