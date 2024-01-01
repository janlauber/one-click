/// <reference path="../pb_data/types.d.ts" />
migrate((db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("zpmsid5xlf5361s")

  // add
  collection.schema.addField(new SchemaField({
    "system": false,
    "id": "tj2ookse",
    "name": "project",
    "type": "relation",
    "required": true,
    "presentable": false,
    "unique": false,
    "options": {
      "collectionId": "7kff2zw80a7rmbu",
      "cascadeDelete": false,
      "minSelect": null,
      "maxSelect": 1,
      "displayFields": null
    }
  }))

  return dao.saveCollection(collection)
}, (db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("zpmsid5xlf5361s")

  // remove
  collection.schema.removeField("tj2ookse")

  return dao.saveCollection(collection)
})
