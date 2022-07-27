export interface IBatchSessionTopics {
  id?: string
  tenantID: string
  topicName: string
  subTopics: IBatchSessionTopics[]
  topicID: string
  moduleID: string
  order: number
  isCompleted: boolean
  totalTime: number
  batchSessionTopicID: string
  batchSessionID: string
  batchSessionTopicOrder: number
  batchSessionTopicIsCompleted: boolean
  batchSessionTopicTotalTime: number
  isChecked: boolean
}

export interface IBatchSessionModules {
  id?: string
  logo: string
  moduleName: string
  moduleTopics: IBatchSessionTopics[]
  resources: string
}




// "id": "c9cf575e-4b89-47c5-bc97-afd2376b67cd",
// "createdAt": "2022-01-27T05:32:30Z",
// "tenantID": "7ca2664b-f379-43db-bdf9-7fdd40707219",
// "topicName": "Coding Standards",
// "subTopics": null,
// "topicID": "44fa345c-8cb4-4d55-a726-68a7cdcc04e5",
// "moduleID": "1ba2d87b-c741-4fe4-9bdd-596d8a55b217",
// "order": 1,
// "isCompleted": null,
// "totalTime": 60,
// "batchSessionTopicID": "00000000-0000-0000-0000-000000000000",
// "batchSessionID": "00000000-0000-0000-0000-000000000000",
// "batchSessionTopicOrder": 0,
// "batchSessionTopicIsCompleted": null,
// "batchSessionTopicTotalTime": 0