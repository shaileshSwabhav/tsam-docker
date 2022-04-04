import { TestBed } from '@angular/core/testing';

import { ProgrammingConceptService } from './programming-concept.service';

describe('ProgrammingConceptService', () => {
  let service: ProgrammingConceptService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(ProgrammingConceptService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
