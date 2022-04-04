import { TestBed } from '@angular/core/testing';

import { CompanyRequirementService } from './company-requirement.service';

describe('CompanyRequirementService', () => {
  let service: CompanyRequirementService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(CompanyRequirementService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
