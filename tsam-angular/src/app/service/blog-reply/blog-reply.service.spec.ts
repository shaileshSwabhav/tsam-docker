import { TestBed } from '@angular/core/testing';

import { BlogReplyService } from './blog-reply.service';

describe('BlogReplyService', () => {
  let service: BlogReplyService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(BlogReplyService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
