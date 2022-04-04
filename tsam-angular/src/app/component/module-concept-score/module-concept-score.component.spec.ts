import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { ModuleConceptScoreComponent } from './module-concept-score.component';

describe('ModuleConceptScoreComponent', () => {
  let component: ModuleConceptScoreComponent;
  let fixture: ComponentFixture<ModuleConceptScoreComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ ModuleConceptScoreComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ModuleConceptScoreComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
