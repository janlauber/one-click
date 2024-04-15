/// <reference path="../pb_data/types.d.ts" />
migrate((db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("zpmsid5xlf5361s")

  // add
  collection.schema.addField(new SchemaField({
    "system": false,
    "id": "xkztc4vt",
    "name": "deployment",
    "type": "relation",
    "required": false,
    "presentable": false,
    "unique": false,
    "options": {
      "collectionId": "h2e1cdq94xgdclh",
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
  collection.schema.removeField("xkztc4vt")

  return dao.saveCollection(collection)
})
