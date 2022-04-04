import { TestBed } from '@angular/core/testing';

import { EnquiryFormService } from './enquiry-form.service';

describe('EnquiryFormService', () => {
  let service: EnquiryFormService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(EnquiryFormService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
