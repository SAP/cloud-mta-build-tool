rmdir moduleNew\symlink_dir_to_content

rmdir content\symlink_dir_to_another_content

del another_content\symlink_to_test4.txt

rmdir symlink_dir_to_moduleNew

del symlink_to_symlink_broken

rmdir symlink_dir_to_symlink_dir_broken

rmdir dir_with_recursive_symlink\subdir\symlink_dir_recursion_to_parent_dir

rmdir dir_with_recursive_symlink\subdir\symlink_dir_to_sibling

rmdir dir_with_recursive_symlink\subdir2\symlink_dir_to_sibling