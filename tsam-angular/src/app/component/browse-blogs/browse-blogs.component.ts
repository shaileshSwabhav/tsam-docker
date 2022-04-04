import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { BlogService, IBlogSnippetDTO } from 'src/app/service/blog/blog.service';
import { GeneralService } from 'src/app/service/general/general.service';
import { UtilityService } from 'src/app/service/utility/utility.service';
import { SpinnerService } from 'src/app/service/spinner/spinner.service';;

@Component({
  selector: 'app-browse-blogs',
  templateUrl: './browse-blogs.component.html',
  styleUrls: ['./browse-blogs.component.css']
})
export class BrowseBlogsComponent implements OnInit {

  // Components.
  blogTopicList: any[]
  selectedBlogTopicList: any[]

  // Blog snippet.
  blogSnippetList: IBlogSnippetDTO[]

  // Image.
  defaultProfileImage: string

  // Spinner.



  // Pagination.
  limitForLatest: number
  limitForTrending: number
  currentPageForLatest: number
  currentPageForTrending: number
  totalLatestBlogSnippets: number
  totalTrendingBlogSnippets: number
  offsetForLatest: number
  offsetForTrending: number
  paginationStart: number
  paginationEnd: number

  // Tab.
  selectedTab: number

  // Banner image.
  defaultBannerImage: string

  constructor(
    private router: Router,
    private generalService: GeneralService,
    private utilityService: UtilityService,
    private spinnerService: SpinnerService,
    private blogService: BlogService,
  ) {
    this.initializeVariables()
    this.getBlogTopicList()
    this.getAllTrendingBlogSnippet()
  }


  get ongoingOperations() {
    return this.spinnerService.ongoingOperations
  }

  ngOnInit(): void { }

  // Initialize global variables.
  initializeVariables(): void {

    // Components.
    this.blogTopicList = []
    this.selectedBlogTopicList = []

    // Image.
    this.defaultProfileImage = "assets/images/default-profile-image.png"

    // Spinner.
    this.spinnerService.loadingMessage = "Getting Blogs"


    // Pagination.
    this.limitForLatest = 5
    this.offsetForLatest = 0
    this.currentPageForLatest = 0
    this.limitForTrending = 5
    this.offsetForTrending = 0
    this.currentPageForTrending = 0
    this.paginationStart = 0
    this.paginationEnd = 0

    // Tab.
    this.selectedTab = 1

    // Banner image.
    this.defaultBannerImage = "assets/images/blog-banner-image-default.jpg"
  }

  // ================================================OTHER FUNCTIONS FOR BLOG===============================================

  // Page change for latest blogs.
  changePageForLatestBlogs(pageNumber: number): void {
    this.currentPageForLatest = pageNumber
    this.offsetForLatest = this.currentPageForLatest - 1
    this.getAllLatestBlogSnippet()
  }

  // Page change for trending blogs.
  changePageForTrendingBlogs(pageNumber: number): void {
    this.currentPageForTrending = pageNumber
    this.offsetForTrending = this.currentPageForTrending - 1
    this.getAllTrendingBlogSnippet()
  }

  // Set total list on current page.
  setPaginationString(limit: number, offset: number, total: number): void {
    this.paginationStart = limit * offset + 1
    this.paginationEnd = +limit + limit * offset
    if (total < this.paginationEnd) {
      this.paginationEnd = total
    }
  }

  // Format fields of blog snippet list.
  formatBlogSnippetListFields(): void {

    for (let i = 0; i < this.blogSnippetList.length; i++) {

      // Give default profile image to authors of blog snippets.
      if (!this.blogSnippetList[i].author.image) {
        this.blogSnippetList[i].author.image = this.defaultProfileImage
      }

      // Give default banner image to all blogs 
      if (!this.blogSnippetList[i].bannerImage) {
        this.blogSnippetList[i].bannerImage = this.defaultBannerImage
      }

      console.log(this.blogSnippetList[i].bannerImage)
    }
  }

  // Redirect to blog details by blog id.
  redirectToBlogDetails(blogID: string): void {
    this.router.navigate(['/community/blog/browse-topics/blog-details'], {
      queryParams: {
        "blogID": blogID,
      }
    }).catch(err => {
      console.error(err)
    })
  }

  // On changing tab get the component lists.
  onTabChange(event: any) {
    this.selectedTab = event
    if (this.selectedTab == 1) {
      this.getAllTrendingBlogSnippet()
    }
    if (this.selectedTab == 2) {
      this.getAllLatestBlogSnippet()
    }
  }

  // On clicking blog topic badge.
  onBlogTopicClick(blogTopic: any, isBlogTopicClick: boolean): void {

    // If blog topic is clicked (right panel)
    if (isBlogTopicClick) {

      // Remove the blog topic from blog topic list.
      for (let i = 0; i < this.blogTopicList.length; i++) {
        if (this.blogTopicList[i].id == blogTopic.id) {
          this.blogTopicList.splice(i, 1)
          break
        }
      }

      // Add the blog topic to selected blog topic list.
      this.selectedBlogTopicList.push(blogTopic)
      if (this.selectedTab == 1) {
        this.getAllTrendingBlogSnippet()
        return
      }
      this.getAllLatestBlogSnippet()
      return
    }

    // If selected blog topic is clicked (left panel)
    // Remove the blog topic from selected blog topic list.
    for (let i = 0; i < this.selectedBlogTopicList.length; i++) {
      if (this.selectedBlogTopicList[i].id == blogTopic.id) {
        this.selectedBlogTopicList.splice(i, 1)
        break
      }
    }

    // Add the blog topic to blog topic list.
    this.blogTopicList.push(blogTopic)

    if (this.selectedTab == 1) {
      this.getAllTrendingBlogSnippet()
      return
    }
    this.getAllLatestBlogSnippet()
  }

  // =============================================================GET FUNCTIONS==========================================================================

  // Get all latest blog snippets.
  getAllLatestBlogSnippet() {
    this.spinnerService.loadingMessage = "Getting All Latest Blogs"


    let queryParams: any = {}
    if (this.selectedBlogTopicList?.length > 0) {
      let blogTopicIDs: string[] = []
      for (let i = 0; i < this.selectedBlogTopicList.length; i++) {
        blogTopicIDs.push(this.selectedBlogTopicList[i].id)
      }
      queryParams.blogTopicIDs = blogTopicIDs
    }
    this.blogService.getAllLatestBlogSnippet(this.limitForLatest, this.offsetForLatest, queryParams).subscribe((response) => {
      this.totalLatestBlogSnippets = response.headers.get('X-Total-Count')
      this.blogSnippetList = response.body
      this.formatBlogSnippetListFields()
    }, error => {
      console.error(error)
      if (error.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
      }
    }).add(() => {
      this.setPaginationString(this.limitForLatest, this.offsetForLatest, this.totalLatestBlogSnippets)
    })
  }

  // Get all trending blog snippets.
  getAllTrendingBlogSnippet() {
    this.spinnerService.loadingMessage = "Getting All Trending Blogs"


    let queryParams: any = {}
    if (this.selectedBlogTopicList?.length > 0) {
      let blogTopicIDs: string[] = []
      for (let i = 0; i < this.selectedBlogTopicList.length; i++) {
        blogTopicIDs.push(this.selectedBlogTopicList[i].id)
      }
      queryParams.blogTopicIDs = blogTopicIDs
    }
    this.blogService.getAllTrendingBlogSnippet(this.limitForTrending, this.offsetForTrending, queryParams).subscribe((response) => {
      this.totalTrendingBlogSnippets = response.headers.get('X-Total-Count')
      this.blogSnippetList = response.body
      this.formatBlogSnippetListFields()
    }, error => {
      console.error(error)
      if (error.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
      }
    }).add(() => {
      this.setPaginationString(this.limitForTrending, this.offsetForTrending, this.totalTrendingBlogSnippets)
    })
  }

  // Get blog topic list.
  getBlogTopicList(): void {
    this.generalService.getBlogTopicList().subscribe(response => {
      this.blogTopicList = response
    }, err => {
      console.error(this.utilityService.getErrorString(err))
    })
  }

}
