import { TestBed } from '@angular/core/testing';

import { BlogReactionService } from './blog-reaction.service';

describe('BlogReactionService', () => {
  let service: BlogReactionService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(BlogReactionService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
