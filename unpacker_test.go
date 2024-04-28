package unpacker

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

var loadedTestInput, loadedTestExpected string

func TestMain(m *testing.M) {
	loadTestData := func(path string) string {
		data, err := os.ReadFile(path)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Test setup failed:", err)
			os.Exit(1)
		}
		return string(data)
	}
	loadedTestInput = loadTestData(filepath.FromSlash("testdata/test-packer-62-input.js"))
	loadedTestExpected = loadTestData(filepath.FromSlash("testdata/test-packer-non62-input.js"))
	os.Exit(m.Run())
}

func TestNewDEUnpacker(t *testing.T) {
	test := func(input string, expected bool) {
		if _, res := NewDEUnpacker(input); res != expected {
			t.Fatalf("Input:\n%s\n\nValid: %v != %v (expected)\n", input, res, expected)
		}
	}
	test("", false)
	test("var a = b", false)
	test("eval(function(p,a,c,k,e,r", true)
	test("eval ( function(p, a, c, k, e, r", true)
}

func TestUnpack(t *testing.T) {
	test := func(input, expected string) {
		unpacker, valid := NewDEUnpacker(input)
		if !valid {
			t.Fatalf("Input:\n%s\n\nInput should have been valid, but wasn't.\n", input)
		}
		res, err := unpacker.Unpack()
		if err != nil {
			t.Fatalf("Input:\n%s\n\nAn error occurred while unpacking: %s\n", input, err)
		}
		if res != expected {
			t.Fatalf("Input:\n%s\n\nOutput:\n%s\n\nExpected:\n%s\n", input, res, expected)
		}
	}
	test(
		"eval(function(p,a,c,k,e,r){e=String;if(!''.replace(/^/,String)){while(c--)r[c]=k[c]||c;k=[function(e){return r[e]}];e=function(){return'\\\\w+'};c=1};while(c--)if(k[c])p=p.replace(new RegExp('\\\\b'+e(c)+'\\\\b','g'),k[c]);return p}('0 2=1',62,3,'var||a'.split('|'),0,{}))",
		"var a=1",
	)
	test(
		"function test (){alert ('This is a test!')}; eval(function(p,a,c,k,e,r){e=String;if(!''.replace(/^/,String)){while(c--)r[c]=k[c]||c;k=[function(e){return r[e]}];e=function(){return'\\w+'};c=1};while(c--)if(k[c])p=p.replace(new RegExp('\\b'+e(c)+'\\b','g'),k[c]);return p}('0 2=\\'{Íâ–+›ï;ã†Ù¥#\\'',3,3,'var||a'.split('|'),0,{}))",
		"function test (){alert ('This is a test!')}; var a='{Íâ–+›ï;ã†Ù¥#'",
	)
	test(
		"eval(function(p,a,c,k,e,d){e=function(c){return c.toString(36)};if(!''.replace(/^/,String)){while(c--){d[c.toString(a)]=k[c]||c.toString(a)}k=[function(e){return d[e]}];e=function(){return'\\w+'};c=1};while(c--){if(k[c]){p=p.replace(new RegExp('\\b'+e(c)+'\\b','g'),k[c])}}return p}('2 0=\"4 3!\";2 1=0.5(/b/6);a.9(\"8\").7=1;',12,12,'str|n|var|W3Schools|Visit|search|i|innerHTML|demo|getElementById|document|w3Schools'.split('|'),0,{}))",
		`var str="Visit W3Schools!";var n=str.search(/w3Schools/i);document.getElementById("demo").innerHTML=n;`,
	)
	test(
		"a=b;\r\nwhile(1){\ng=h;{return'\\w+'};break;eval(function(p,a,c,k,e,d){e=function(c){return c.toString(36)};if(!''.replace(/^/,String)){while(c--){d[c.toString(a)]=k[c]||c.toString(a)}k=[function(e){return d[e]}];e=function(){return'\\w+'};c=1};while(c--){if(k[c]){p=p.replace(new RegExp('\\b'+e(c)+'\\b','g'),k[c])}}return p}('$(5).4(3(){$('.1').0(2);$('.6').0(d);$('.7').0(b);$('.a').0(8);$('.9').0(c)});',14,14,'html|r5e57|8080|function|ready|document|r1655|rc15b|8888|r39b0|r6ae9|3128|65309|80'.split('|'),0,{}))c=abx;",
		"a=b;\r\nwhile(1){\ng=h;{return'\\w+'};break;$(document).ready(function(){$('.r5e57').html(8080);$('.r1655').html(80);$('.rc15b').html(3128);$('.r6ae9').html(8888);$('.r39b0').html(65309)});c=abx;",
	)
	test(
		"eval(function(p,a,c,k,e,r){e=function(c){return c.toString(36)};if('0'.replace(0,e)==0){while(c--)r[e(c)]=k[c];k=[function(e){return r[e]||e}];e=function(){return'[0-9ab]'};c=1};while(c--)if(k[c])p=p.replace(new RegExp('\\b'+e(c)+'\\b','g'),k[c]);return p}('$(5).a(6(){ $('.8').0(1); $('.b').0(4); $('.9').0(2); $('.7').0(3)})',[],12,'html|52136|555|65103|8088|document|function|r542c|r8ce6|rb0de|ready|rfab0'.split('|'),0,{}))",
		"$(document).ready(function(){ $('.r8ce6').html(52136); $('.rfab0').html(8088); $('.rb0de').html(555); $('.r542c').html(65103)})",
	)
	test(
		loadedTestInput,
		loadedTestExpected,
	)
}
