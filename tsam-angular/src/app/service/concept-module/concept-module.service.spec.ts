import { TestBed } from '@angular/core/testing';

import { ConceptModuleService } from './concept-module.service';

describe('ConceptModuleService', () => {
  let service: ConceptModuleService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(ConceptModuleService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
