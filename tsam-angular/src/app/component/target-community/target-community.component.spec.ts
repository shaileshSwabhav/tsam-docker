import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { TargetCommunityComponent } from './target-community.component';

describe('TargetCommunityComponent', () => {
  let component: TargetCommunityComponent;
  let fixture: ComponentFixture<TargetCommunityComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ TargetCommunityComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(TargetCommunityComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
