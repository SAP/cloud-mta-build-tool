
set "PREVIOUS_DIR=%CD%"

cd moduleNew
mklink /D symlink_dir_to_content "..\content"
cd "%PREVIOUS_DIR%"

cd content
mklink /D symlink_dir_to_another_content "..\another_content"
cd "%PREVIOUS_DIR%"

cd another_content
mklink symlink_to_test4.txt "..\test4.txt"
cd "%PREVIOUS_DIR%"

mklink /D symlink_dir_to_moduleNew "moduleNew"

mklink symlink_to_symlink_broken "link_to_broken_symlink"

mklink /D symlink_dir_to_symlink_dir_broken "link_to_broken_symlink_dir"

cd "dir_with_recursive_symlink\subdir"
mklink /D symlink_dir_recursion_to_parent_dir "..\"
mklink /D symlink_dir_to_sibling "..\subdir2\symlink_dir_to_sibling"
cd "%PREVIOUS_DIR%"

cd "dir_with_recursive_symlink\subdir2"
mklink /D symlink_dir_to_sibling "..\subdir\symlink_dir_to_sibling"
cd "%PREVIOUS_DIR%"