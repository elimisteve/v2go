fn main() {
	reLinks := regexp.MustCompile('<a href="(http.*?)" class="storylink"')

	resp := http.Get('https://news.ycombinator.com')?
	defer resp.Body.Close()

	html := ioutil.ReadAll(resp.Body)?

	for link in reLinks.FindAllSubmatch(html, -1) {
		println('${string(link[1])}')
	}
}
