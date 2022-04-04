import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { ModuleConceptScoreNodeComponent } from './module-concept-score-node.component';

describe('ModuleConceptScoreNodeComponent', () => {
  let component: ModuleConceptScoreNodeComponent;
  let fixture: ComponentFixture<ModuleConceptScoreNodeComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ ModuleConceptScoreNodeComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ModuleConceptScoreNodeComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
