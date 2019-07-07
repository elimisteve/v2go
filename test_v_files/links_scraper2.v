fn main() {
	reLinks := regexp.MustCompile('<a href="(http.*?)" class="storylink"')

	url := 'https://news.ycombinator.com'
	resp := http.Get('$url')?

	defer resp.Body.Close()

	html := ioutil.ReadAll(resp.Body)?

	for link in reLinks.FindAllSubmatch(html, -1) {
		println('${string(link[1])}')
	}
}
