import { TestBed } from '@angular/core/testing';

import { BlogViewerService } from './blog-viewer.service';

describe('BlogViewerService', () => {
  let service: BlogViewerService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(BlogViewerService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
