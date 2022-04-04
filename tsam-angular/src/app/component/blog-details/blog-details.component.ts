import { Component, OnInit, ViewChild } from '@angular/core';
import { FormGroup, FormBuilder, FormControl, Validators } from '@angular/forms';
import { ActivatedRoute, Router } from '@angular/router';
import { BlogService, IBlogDTO } from 'src/app/service/blog/blog.service';
import { SpinnerService } from 'src/app/service/spinner/spinner.service';;
import { GeneralService } from 'src/app/service/general/general.service';
import { UtilityService } from 'src/app/service/utility/utility.service';
import { BlogReplyService, IBlogReplyDTO } from 'src/app/service/blog-reply/blog-reply.service';
import { LocalService } from 'src/app/service/storage/local.service';
import { BlogReactionService, IBlogReaction, IBlogReactionDTO } from 'src/app/service/blog-reaction/blog-reaction.service';
import { BlogViewerService, IBlogView } from 'src/app/service/blog-viewer/blog-viewer.service';
import { NgbModalRef, NgbModalOptions, NgbModal } from '@ng-bootstrap/ng-bootstrap';

@Component({
  selector: 'app-blog-details',
  templateUrl: './blog-details.component.html',
  styleUrls: ['./blog-details.component.css']
})
export class BlogDetailsComponent implements OnInit {

  // Components.
  blogTopicList: any[]

  // Blog.
  blog: IBlogDTO
  blogID: string

  // Reply.
  replyToBlogForm: FormGroup
  replyToReplyForm: FormGroup
  replyList: IBlogReplyDTO[]
  selectedReplyID: string
  selectedReplyToBlogID: string

  // Reaction.
  selectedReactionList: IBlogReactionDTO[]

  // Spinner.



  // Image.
  defaultProfileImage: string

  // Login details.
  credentialID: string

  // Flags.
  isOperationReplyToBlogUpdate: boolean
  isOperationReplyToReplyUpdate: boolean

  // Modal.
  modalRef: any
  @ViewChild('showClapsModal') showClapsModal: any

  constructor(
    private formBuilder: FormBuilder,
    private activatedRoute: ActivatedRoute,
    private spinnerService: SpinnerService,
    private blogService: BlogService,
    private blogReplyService: BlogReplyService,
    private blogReactionService: BlogReactionService,
    private blogViewerService: BlogViewerService,
    private generalService: GeneralService,
    private utilityService: UtilityService,
    private localService: LocalService,
    private modalService: NgbModal,
    private router: Router,
  ) {
    this.initializeVariables()
    this.getBlogTopicList()
    this.createReplyToBlogForm()
    this.createReplyToReplyForm()
  }


  get ongoingOperations() {
    return this.spinnerService.ongoingOperations
  }

  ngOnInit(): void { }

  // Initialize global variables.
  initializeVariables(): void {

    // Components.
    this.blogTopicList = []

    // Reply.
    this.replyList = []

    // Spinner.
    this.spinnerService.loadingMessage = "Getting Blog Details"


    // Image.
    this.defaultProfileImage = "assets/images/default-profile-image.png"

    // Login details.
    this.credentialID = this.localService.getJsonValue("credentialID")

    // Flags.
    this.isOperationReplyToBlogUpdate = false
    this.isOperationReplyToReplyUpdate = false

    // Reaction.
    this.selectedReactionList = []

    // Get query params.
    this.blogID = this.activatedRoute.snapshot.queryParamMap.get("blogID")
    if (this.blogID) {
      this.getRepliesByBlogID()
      this.addBlogView()
    }
  }

  // ================================================ CREATE FORM FUNCTIONS ===============================================

  // Create reply to blog form.
  createReplyToBlogForm(): void {
    this.replyToBlogForm = this.formBuilder.group({
      id: new FormControl(null),
      blogID: new FormControl(null),
      replyID: new FormControl(null),
      replierID: new FormControl(null),
      reply: new FormControl(null, [Validators.maxLength(2000)]),
      isBlogReply: new FormControl(true),
      isVerified: new FormControl(true),
    })
  }

  // Create reply to reply form.
  createReplyToReplyForm(): void {
    this.replyToReplyForm = this.formBuilder.group({
      id: new FormControl(null),
      blogID: new FormControl(null),
      replyID: new FormControl(null),
      replierID: new FormControl(null),
      reply: new FormControl(null, [Validators.maxLength(2000)]),
      isBlogReply: new FormControl(false),
      isVerified: new FormControl(true),
    })
  }

  // ================================================ CRUD FUNCTIONS FOR BLOG REPLY ===============================================

  // Add blog reply.
  addBlogReply(replyForm: any): void {
    this.spinnerService.loadingMessage = "Sending Reply"


    this.blogReplyService.addBlogReply(replyForm.value).subscribe((response: any) => {
      this.getRepliesByBlogID()
      alert("Reply sent successsfully")
    }, (error) => {
      alert("Reply could not be sent, please try again later")
      console.error(error)
    })
  }

  // Update blog reply.
  updateBlogReply(replyForm: any): void {
    this.spinnerService.loadingMessage = "Updating Reply"


    this.blogReplyService.updateBlogReply(replyForm.value).subscribe((response: any) => {
      this.getRepliesByBlogID()

      // Make all replies' update form invisible.
      for (let i = 0; i < this.replyList.length; i++) {
        this.replyList[i].isReplyUpdateFormVisible = false
      }

      alert("Reply updated successsfully")
    }, (error) => {
      alert("Reply could not be updated, please try again later")
      console.error(error)
    })
  }

  // Delete blog reply.
  deleteBlogReply(replyID: string): void {
    this.spinnerService.loadingMessage = "Deleting Reply"


    this.blogReplyService.deleteBlogReply(replyID).subscribe((response: any) => {
      this.getRepliesByBlogID()
      alert("Reply deleted successsfully")
    }, (error) => {
      alert("Reply could not be deleted, please try again later")
      console.error(error)
    })
  }

  // ================================================ CRUD FUNCTIONS FOR BLOG REACTION ===============================================

  // Add blog reaction.
  addBlogReaction(blogOrReply: any, blogReaction: any, isClap: boolean): void {
    this.spinnerService.loadingMessage = "Sending Reaction"


    this.blogReactionService.addBlogReaction(blogReaction).subscribe((response: any) => {
      blogOrReply.isLoggedInClap = isClap
      this.getBlogByID()
      this.getRepliesByBlogID()
    }, (error) => {
      alert("Reaction could not be sent, please try again later")
      console.error(error)
    })
  }

  // Update blog reaction.
  updateBlogReaction(blogOrReply: any, blogReaction: any, isClap: boolean): void {
    this.spinnerService.loadingMessage = "Sending Reaction"


    blogReaction.reactorID = blogReaction.reactor.id
    delete blogReaction.reactor
    this.blogReactionService.updateBlogReaction(blogReaction).subscribe((response: any) => {
      blogOrReply.isLoggedInClap = isClap
      this.getBlogByID()
      this.getRepliesByBlogID()
    }, (error) => {
      alert("Reaction could not be sent, please try again later")
      console.error(error)
    })
  }

  // Delete blog reaction.
  deleteBlogReaction(blogOrReply: any, reactionID: string): void {
    this.spinnerService.loadingMessage = "Deleting Reaction"


    this.blogReactionService.deleteBlogReaction(reactionID).subscribe((response: any) => {
      blogOrReply.isLoggedInClap = null
      this.getBlogByID()
      this.getRepliesByBlogID()
    }, (error) => {
      alert("Reaction could not be deleted, please try again later")
      console.error(error)
    })
  }

  // ================================================ ON BUTTON CLICK FUNCTIONS ===============================================

  // On clicking reply submit button.
  onReplySubmit(replyForm: any): void {

    // If reply is blank then dont submit.
    if (replyForm.get('reply').value == null || replyForm.get('reply').value == "") {
      return
    }

    // If reply form is reply to blog from.
    if (replyForm.get('isBlogReply').value) {
      if (this.isOperationReplyToBlogUpdate) {
        this.updateBlogReply(replyForm)
        this.createReplyToBlogForm()
        return
      }
      replyForm.get('blogID').setValue(this.blogID)
      replyForm.get('replierID').setValue(this.credentialID)
      this.addBlogReply(replyForm)
      this.createReplyToBlogForm()
      return
    }

    // If reply form is reply to reply from.
    if (this.isOperationReplyToReplyUpdate) {
      this.updateBlogReply(replyForm)
      this.createReplyToReplyForm()
      return
    }
    replyForm.get('replyID').setValue(this.selectedReplyID)
    replyForm.get('replierID').setValue(this.credentialID)
    this.addBlogReply(replyForm)
    this.createReplyToReplyForm()
  }

  // On clicking show all replies text.
  onShowRepliesClick(reply: IBlogReplyDTO): void {
    if (reply.isRepliesVisible) {
      reply.isRepliesVisible = false
      return
    }
    reply.isRepliesVisible = true
  }

  // On clicking reply text.
  onReplyButtonClick(replyID: string): void {

    this.createReplyToReplyForm()

    this.isOperationReplyToReplyUpdate = false

    // Make the sub replies' update reply form invisible.
    for (let i = 0; i < this.replyList.length; i++) {
      for (let j = 0; j < this.replyList[i].replies.length; j++) {
        this.replyList[i].replies[j].isReplyUpdateFormVisible == false
      }
    }

    // If selected reply's reply text is clicked then hide the reply form.
    if (this.selectedReplyID == replyID) {
      for (let i = 0; i < this.replyList.length; i++) {
        if (this.replyList[i].id == this.selectedReplyID) {
          this.replyList[i].isReplyToReplyFormVisible = false
          this.selectedReplyID = null
          return
        }
      }
    }

    // If selected reply is not present then make the reply as the selected reply.
    if (!this.selectedReplyID) {
      this.selectedReplyID = replyID
    }

    // Make the selected reply's reply form visible and hide other replies' reply form.    
    for (let i = 0; i < this.replyList.length; i++) {
      if (this.replyList[i].id == replyID) {
        this.selectedReplyID = this.replyList[i].id
        this.replyList[i].isReplyToReplyFormVisible = true
        continue
      }
      this.replyList[i].isReplyToReplyFormVisible = false
    }
  }

  // On clicking update reply button.
  onUpdateReplyButtonClick(reply: IBlogReplyDTO): void {

    // ************************** If reply to blog is updated **************************.

    if (reply.isBlogReply) {

      this.isOperationReplyToBlogUpdate = true

      // Make the selected reply's update form visible and hide other replies' update form.
      for (let i = 0; i < this.replyList.length; i++) {
        if (this.replyList[i].id == reply.id) {
          this.replyList[i].isReplyUpdateFormVisible = true
          continue
        }
        this.replyList[i].isReplyUpdateFormVisible = false
      }

      this.createReplyToBlogForm()
      this.replyToBlogForm.patchValue(reply)
      this.replyToBlogForm.get('replierID').setValue(reply.replier.id)
      return
    }

    // ************************** If reply to reply is updated **************************.

    this.isOperationReplyToReplyUpdate = true

    // Make the selected reply's reply form visible and hide other replies' reply form.
    for (let i = 0; i < this.replyList.length; i++) {
      if (this.replyList[i].id == reply.replyID) {
        this.selectedReplyID = this.replyList[i].id
        this.replyList[i].isReplyToReplyFormVisible = true
        continue
      }
      this.replyList[i].isReplyToReplyFormVisible = false
    }

    // Make the selected reply's reply update form visible and hide other replies' reply update form.
    for (let i = 0; i < this.replyList.length; i++) {
      for (let j = 0; j < this.replyList[i].replies.length; j++) {
        if (this.replyList[i].replies[j].id == reply.id) {
          this.replyList[i].replies[j].isReplyUpdateFormVisible == true
          continue
        }
        this.replyList[i].replies[j].isReplyUpdateFormVisible == false
      }
    }

    this.createReplyToReplyForm()
    this.replyToReplyForm.patchValue(reply)
    this.replyToReplyForm.get('replierID').setValue(reply.replier.id)
  }

  // On clicking delete reply button.
  onDeleteReplyButtonClick(replyID: string): void {
    if (confirm("Are you sure you want to delete this reply")) {
      this.deleteBlogReply(replyID)
    }
  }

  // On clicking cancel reply button.
  onCancelReplyButtonClick(isBlogReply: boolean): void {

    // ************************** If reply to blog is update is cancelled **************************.

    if (isBlogReply) {
      this.isOperationReplyToBlogUpdate = false

      // Make all replies' update form invisible.
      for (let i = 0; i < this.replyList.length; i++) {
        this.replyList[i].isReplyUpdateFormVisible = false
      }

      this.createReplyToBlogForm()

      return
    }

    // ************************** If reply to reply is update is cancelled **************************.

    this.isOperationReplyToReplyUpdate = false

    // Make all replies' replies' update form invisible.
    for (let i = 0; i < this.replyList.length; i++) {
      for (let j = 0; j < this.replyList[i].replies.length; j++) {
        this.replyList[i].replies[j].isReplyUpdateFormVisible == false
      }
    }

    // Make all replies' reply to reply form invisible.
    for (let i = 0; i < this.replyList.length; i++) {
      this.replyList[i].isReplyToReplyFormVisible = false
    }

    this.createReplyToReplyForm()
  }

  // On clicking blog reaction icon.
  onBlogReactionClick(blogOrReply: any, isClap: boolean, isBolg: boolean): void {

    // If there was no reaction then add reaction.
    if (blogOrReply.isLoggedInClap == null) {
      let blogReaction: IBlogReaction = {
        reactorID: this.credentialID,
        isClap: isClap
      }
      if (isBolg) {
        blogReaction.blogID = blogOrReply.id
      }
      else {
        blogReaction.replyID = blogOrReply.id
      }
      this.addBlogReaction(blogOrReply, blogReaction, isClap)
      return
    }

    // If isLoggedIn reaction is clicked again then delete reaction.
    if (blogOrReply.isLoggedInClap == isClap) {
      for (let i = 0; i < blogOrReply.reactions.length; i++) {
        if (blogOrReply.reactions[i].reactor.id == this.credentialID) {
          this.deleteBlogReaction(blogOrReply, blogOrReply.reactions[i].id)
          return
        }
      }
    }

    // If isLoggedIn reaction's opposite reaction is clicked then update reaction.
    if (blogOrReply.isLoggedInClap != isClap) {
      for (let i = 0; i < blogOrReply.reactions.length; i++) {
        if (blogOrReply.reactions[i].reactor.id == this.credentialID) {
          blogOrReply.reactions[i].isClap = isClap
          this.updateBlogReaction(blogOrReply, blogOrReply.reactions[i], isClap)
          return
        }
      }
    }
  }

  // On clicking claps count number.
  onClapsCountClick(replyReactions?: IBlogReactionDTO[]): void {
    this.selectedReactionList = []
    if (this.blog.author.id == this.credentialID) {
      if (replyReactions && replyReactions?.length > 0) {
        this.selectedReactionList = replyReactions
      }
      else if (this.blog.clapReactions?.length > 0) {
        this.selectedReactionList = this.blog.clapReactions
      }
      if (this.selectedReactionList.length > 0) {
        this.formatReactionFieldListFields()
        this.openModal(this.showClapsModal, 'lg')
      }
    }
  }

  // ================================================ FORMAT FUNCTIONS FOR BLOG ===============================================

  // Format fields of blog.
  formatBlogFields(): void {

    // Give default profile image to author of blog.
    if (!this.blog.author.image) {
      this.blog.author.image = this.defaultProfileImage
    }

    // Give default banner image to blog.
    if (!this.blog.bannerImage) {
      this.blog.bannerImage = "assets/images/blog-banner-image-default.jpg"
    }

    // Set blog reactions.
    this.blog.isLoggedInClap = null
    this.blog.clapReactions = []
    this.blog.slapReactions = []
    for (let i = 0; i < this.blog.reactions?.length; i++) {

      // If author has reacted on the blog then set the isClap field of blog.
      if (this.blog.reactions[i].reactor.id == this.credentialID) {
        this.blog.isLoggedInClap = this.blog.reactions[i].isClap
      }

      // Get clap and slap count.
      if (this.blog.reactions[i].isClap) {
        this.blog.clapReactions.push(this.blog.reactions[i])
        continue
      }
      this.blog.slapReactions.push(this.blog.reactions[i])
    }
  }

  // Format fields of reply list.
  formatReplyFields(): void {

    for (let i = 0; i < this.replyList.length; i++) {

      // Set isRepliesVisible, isReplyFormVisible and isBlogReply as false.
      this.replyList[i].isRepliesVisible = false
      this.replyList[i].isReplyToReplyFormVisible = false
      this.replyList[i].isReplyUpdateFormVisible = false
      this.replyList[i].isBlogReply = true

      // Give default profile image to repliers of reply list.
      if (!this.replyList[i].replier.image) {
        this.replyList[i].replier.image = "assets/images/default-profile-image.png"
      }

      // Set reply reactions.
      this.replyList[i].isLoggedInClap = null
      this.replyList[i].clapReactions = []
      this.replyList[i].slapReactions = []
      for (let j = 0; j < this.replyList[i].reactions?.length; j++) {

        // If author has reacted on the reply then set the isClap field of reply.
        if (this.replyList[i].reactions[j].reactor.id == this.credentialID) {
          this.replyList[i].isLoggedInClap = this.replyList[i].reactions[j].isClap
        }

        // Get clap and slap count.
        if (this.replyList[i].reactions[j].isClap) {
          this.replyList[i].clapReactions.push(this.replyList[i].reactions[j])
          continue
        }
        this.replyList[i].slapReactions.push(this.replyList[i].reactions[j])
      }

      for (let j = 0; j < this.replyList[i].replies?.length; j++) {

        // Set reply reactions.
        this.replyList[i].replies[j].isLoggedInClap = null
        this.replyList[i].replies[j].clapReactions = []
        this.replyList[i].replies[j].slapReactions = []
        for (let k = 0; k < this.replyList[i].replies[j]?.reactions?.length; k++) {

          // If author has reacted on the reply then set the isClap field of reply.
          if (this.replyList[i].replies[j].reactions[k].reactor.id == this.credentialID) {
            this.replyList[i].replies[j].isLoggedInClap = this.replyList[i].replies[j].reactions[k].isClap
          }

          // Get clap and slap count.
          if (this.replyList[i].replies[j].reactions[k].isClap) {
            this.replyList[i].replies[j].clapReactions.push(this.replyList[i].replies[j].reactions[k])
            continue
          }
          this.replyList[i].replies[j].slapReactions.push(this.replyList[i].replies[j].reactions[k])
        }

        // Give default profile image to repliers of reply list.
        if (!this.replyList[i].replies[j].replier.image) {
          this.replyList[i].replies[j].replier.image = "assets/images/default-profile-image.png"
        }

        // Set isBlogReply as false.
        this.replyList[i].replies[j].isBlogReply = false
      }
    }
  }

  // Format fields of selected reaction list.
  formatReactionFieldListFields(): void {

    // Give default profile image to repliers of reply list.
    for (let i = 0; i < this.selectedReactionList.length; i++) {
      if (!this.selectedReactionList[i].reactor.image) {
        this.selectedReactionList[i].reactor.image = "assets/images/default-profile-image.png"
      }
    }
  }

  // ================================================ OTHER FUNCTIONS FOR BLOG ===============================================

  // Used to open modal.
  openModal(content: any, size?: string): NgbModalRef {
    if (!size) {
      size = 'sm'
    }
    let options: NgbModalOptions = {
      ariaLabelledBy: 'modal-basic-title', keyboard: false,
      backdrop: 'static', size: size
    }
    this.modalRef = this.modalService.open(content, options)
    return this.modalRef
  }

  // Redirect to browse blogs.
  redirectToBrowseBlogs(): void {
    this.router.navigate(['/community/blog/browse-topics'], {
    }).catch(err => {
      console.error(err)
    })
  }

  // =============================================================GET FUNCTIONS==========================================================================

  // Get blog by id.
  getBlogByID() {
    this.spinnerService.loadingMessage = "Getting All Blogs"


    this.blogService.getBlogByID(this.blogID).subscribe((response) => {
      this.blog = response
      this.formatBlogFields()
    }, error => {
      console.error(error)
      if (error.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
      }
    })
  }

  // Get all replies by blog id.
  getRepliesByBlogID() {
    this.spinnerService.loadingMessage = "Getting All Replies"


    this.blogReplyService.getAllBlogRelpies(this.blogID).subscribe((response) => {
      this.replyList = response
      this.formatReplyFields()
    }, error => {
      console.error(error)
      if (error.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
      }
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

  // Add blog view.
  addBlogView(): void {
    let blogView: IBlogView = {
      blogID: this.blogID,
      viewerID: this.credentialID
    }
    this.blogViewerService.addBlogView(blogView).subscribe(response => {
      this.getBlogByID()
    }, err => {
      console.error(this.utilityService.getErrorString(err))
    })
  }

}
