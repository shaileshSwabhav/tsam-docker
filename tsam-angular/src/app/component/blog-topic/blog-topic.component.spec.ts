import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { BlogTopicComponent } from './blog-topic.component';

describe('BlogTopicComponent', () => {
  let component: BlogTopicComponent;
  let fixture: ComponentFixture<BlogTopicComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ BlogTopicComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(BlogTopicComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
