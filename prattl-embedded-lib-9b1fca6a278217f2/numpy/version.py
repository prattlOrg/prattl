
"""
Module to expose more detailed version info for the installed `numpy`
"""
version = "2.0.0rc2"
__version__ = version
full_version = version

git_revision = "e09a975364d19a562af056c1184affa7f81b812e"
release = 'dev' not in version and '+' not in version
short_version = version.split("+")[0]
