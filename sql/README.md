The python script populates the mysql database with the Chinese characters contained in
chinese_characters.json.

`chinese_characters.json` is the same file being sent to the client on the browser, 
it's the same file as seen in `vue/variant-WordData.json`

It came from https://github.com/gitqwerty777/Chinese-Characters-Standards

I probably need to trim the extra fields from it.  If somebody wants to write a Python
script to do that and put it in a pull request, I'd appreciate that.