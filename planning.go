package main

/* Planning for how to create a tool to collate web resources based on simple
website => keyword/category rules. */

/*
# Design One: Web scraping

## Components
 - Scraper: Responsible for scraping a website and returning the HTML.
 - Provider: Responsible for providing results, to a given format, for a given website type.
 - Enricher: Responsible for getting additional data about a web resource.
 - Filter: Responsible for filtering out results given filter
 - Configuration: Responsible for holding user provided configuration from disk.
 - Server: Responsible for displaying collated web resources as local web page

## Architecture
				Base Website Scraper
	                    |
                 |──────┴────┐
            Provider #1     Provider #2
                 |              |
     List<Resources>      List<Resources>
                 |              |
         Enricher #1 [1]    pass-through
	             └───────┬──────┴
		        	  Filter  <= Configuration
						 |
                   Web resources -> cache/disk
				         |
				         └-> Server

[1] : We would then maybe need a scraper to get details about the resource itself
*/
//
//// Scraper is responsible for scraping a website and returning the HTML.
//type Scraper struct {
//	urlQueue  chan string
//	htmlQueue chan goquery.Document
//}
//
//// Resource is responsible for holding all information relevant to a web resource (e.g article, repo, paper)
//type Resource struct {
//	uri string
//}
//
//// Enricher is responsible for getting additional data about a web resource.
//type Enricher interface {
//	Enrich(resource *Resource) error // mutates underlying data
//}
//
//// ResourceFilter is responsible for filtering out results given filter
//type ResourceFilter interface {
//	Filter(resource *Resource) bool
//}
//
//// ConfigurationManager is responsible for holding user provided configuration from disk.
//type ConfigurationManager struct {
//}
//
//// Server is responsible for displaying collated web resources as local web page
//type Server struct {
//}

/*
# Design Two:
*/
