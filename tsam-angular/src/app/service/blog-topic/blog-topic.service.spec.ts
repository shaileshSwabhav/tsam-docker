import { TestBed } from '@angular/core/testing';

import { BlogTopicService } from './blog-topic.service';

describe('BlogTopicService', () => {
  let service: BlogTopicService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(BlogTopicService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
