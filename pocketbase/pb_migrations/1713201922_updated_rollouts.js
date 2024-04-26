/// <reference path="../pb_data/types.d.ts" />
migrate((db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("22k6ts6gvnp46mc")

  // add
  collection.schema.addField(new SchemaField({
    "system": false,
    "id": "tbscll1i",
    "name": "deployment",
    "type": "relation",
    "required": true,
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
  const collection = dao.findCollectionByNameOrId("22k6ts6gvnp46mc")

  // remove
  collection.schema.removeField("tbscll1i")

  return dao.saveCollection(collection)
})
